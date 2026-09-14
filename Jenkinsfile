pipeline {
  agent any

  options {
    timestamps()
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '10'))
  }

  environment {
    SERVER_DIR = 'server'
    CLIENT_DIR = 'client'
    SERVER_BINARY = 'journal-api'
    GOOSE_BINARY = 'goose'
    GOOSE_INSTALL_PATH = '/usr/local/bin/goose'
    GOOSE_VERSION = 'v3.27.3'
    DEPLOY_HOST = 'host.docker.internal'
    DEPLOY_USER = 'deploy'
    SERVER_SERVICE = 'journal-api'
    SERVER_INSTALL_DIR = '/opt/journal/server'
    SERVER_ENV_FILE = '/etc/journal/server.env'
    WEB_INSTALL_DIR = '/var/www/job4j.ru/journal'
    RUN_MIGRATIONS = 'true'
    ARTIFACT_DIR = "${WORKSPACE}/.local/jenkins-build"
    GO_VERSION = '1.25.7'
    GO_ROOT = "${WORKSPACE}/.local/go"
    GOROOT = "${GO_ROOT}"
    NODE_VERSION = '24.19.0'
    NODE_ROOT = "${WORKSPACE}/.local/node"
    PNPM_VERSION = '11.19.0'
    PNPM_ROOT = "${WORKSPACE}/.local/pnpm"
    GOCACHE = "${WORKSPACE}/.local/go-build"
    NPM_CONFIG_CACHE = "${WORKSPACE}/.local/npm-cache"
    PNPM_STORE_DIR = "${WORKSPACE}/.local/pnpm-store"
    SSH_CREDENTIALS_ID = 'job4j-pixel-deploy-ssh'
    PATH = "${GO_ROOT}/bin:${NODE_ROOT}/bin:${PNPM_ROOT}/bin:${env.PATH}"
  }

  stages {
    stage('Prepare') {
      steps {
        sh '''
          set -eu
          mkdir -p "$GOCACHE" "$NPM_CONFIG_CACHE" "$PNPM_STORE_DIR" "$ARTIFACT_DIR"

          if [ ! -x "$GO_ROOT/bin/go" ] ||
            ! "$GO_ROOT/bin/go" version | grep -q "go${GO_VERSION}"; then
            GO_ARCHIVE="go${GO_VERSION}.linux-amd64.tar.gz"
            GO_DOWNLOAD_URL="https://go.dev/dl/${GO_ARCHIVE}"
            GO_TMP_DIR="${WORKSPACE}/.local/go-install"

            rm -rf "$GO_TMP_DIR"
            mkdir -p "$GO_TMP_DIR"
            curl -fsSL "$GO_DOWNLOAD_URL" -o "$GO_TMP_DIR/$GO_ARCHIVE"
            tar -C "$GO_TMP_DIR" -xzf "$GO_TMP_DIR/$GO_ARCHIVE"
            rm -rf "$GO_ROOT"
            mv "$GO_TMP_DIR/go" "$GO_ROOT"
            rm -rf "$GO_TMP_DIR"
          fi

          if [ ! -x "$NODE_ROOT/bin/node" ] ||
            ! "$NODE_ROOT/bin/node" --version | grep -q "v${NODE_VERSION}"; then
            NODE_DIST="node-v${NODE_VERSION}-linux-x64"
            NODE_ARCHIVE="${NODE_DIST}.tar.gz"
            NODE_DOWNLOAD_URL="https://nodejs.org/dist/v${NODE_VERSION}/${NODE_ARCHIVE}"
            NODE_TMP_DIR="${WORKSPACE}/.local/node-install"

            rm -rf "$NODE_TMP_DIR"
            mkdir -p "$NODE_TMP_DIR"
            curl -fsSL "$NODE_DOWNLOAD_URL" -o "$NODE_TMP_DIR/$NODE_ARCHIVE"
            tar -C "$NODE_TMP_DIR" -xzf "$NODE_TMP_DIR/$NODE_ARCHIVE"
            rm -rf "$NODE_ROOT"
            mv "$NODE_TMP_DIR/$NODE_DIST" "$NODE_ROOT"
            rm -rf "$NODE_TMP_DIR"
          fi

          if [ ! -x "$PNPM_ROOT/bin/pnpm" ] ||
            ! "$PNPM_ROOT/bin/pnpm" --version | grep -q "^${PNPM_VERSION}$"; then
            rm -rf "$PNPM_ROOT"
            npm install --global --prefix "$PNPM_ROOT" "pnpm@${PNPM_VERSION}"
          fi

          go version
          node --version
          pnpm --version
          GOBIN="$ARTIFACT_DIR" go install \
            "github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}"
        '''
      }
    }

    stage('Server: Test') {
      steps {
        dir(env.SERVER_DIR) {
          sh 'go test ./...'
          sh 'go vet ./...'
        }
      }
    }

    stage('Server: Build') {
      steps {
        dir(env.SERVER_DIR) {
          sh '''
            set -eu
            GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
              -o "${ARTIFACT_DIR}/${SERVER_BINARY}" ./cmd
          '''
        }
      }
    }

    stage('Client: Install') {
      steps {
        dir(env.CLIENT_DIR) {
          sh 'pnpm install --frozen-lockfile --store-dir "$PNPM_STORE_DIR"'
        }
      }
    }

    stage('Client: Check') {
      steps {
        dir(env.CLIENT_DIR) {
          sh 'pnpm test -- --maxWorkers=1'
          sh 'pnpm run lint'
          sh 'VITE_BASE_PATH="/journal/" pnpm run build'
        }
      }
    }

    stage('Deploy Server') {
      when {
        anyOf {
          branch 'main'
          expression { !env.BRANCH_NAME }
        }
      }
      steps {
        sshagent(credentials: [env.SSH_CREDENTIALS_ID]) {
          sh '''
            set -eu
            DEPLOY_TARGET="${DEPLOY_USER}@${DEPLOY_HOST}"
            SSH_OPTS="-o StrictHostKeyChecking=no"
            SSH_OPTS="$SSH_OPTS -o UserKnownHostsFile=/dev/null -o BatchMode=yes"

            tar -C "$SERVER_DIR" -czf "$ARTIFACT_DIR/server-assets.tar.gz" \
              api migrations

            ssh $SSH_OPTS "$DEPLOY_TARGET" \
              "sudo -n mkdir -p '${SERVER_INSTALL_DIR}' &&
               sudo -n chown '${DEPLOY_USER}:${DEPLOY_USER}' \
                 '${SERVER_INSTALL_DIR}' &&
               test -r '${SERVER_ENV_FILE}'"

            scp $SSH_OPTS "$ARTIFACT_DIR/$SERVER_BINARY" \
              "$DEPLOY_TARGET:/tmp/$SERVER_BINARY"
            scp $SSH_OPTS "$ARTIFACT_DIR/$GOOSE_BINARY" \
              "$DEPLOY_TARGET:/tmp/$GOOSE_BINARY"
            scp $SSH_OPTS "$ARTIFACT_DIR/server-assets.tar.gz" \
              "$DEPLOY_TARGET:/tmp/journal-server-assets.tar.gz"

            ssh $SSH_OPTS "$DEPLOY_TARGET" \
              "sudo -n install -m 0755 '/tmp/${SERVER_BINARY}' \
                 '${SERVER_INSTALL_DIR}/${SERVER_BINARY}' &&
               rm -f '/tmp/${SERVER_BINARY}' &&
               sudo -n install -m 0755 '/tmp/${GOOSE_BINARY}' \
                 '${GOOSE_INSTALL_PATH}' &&
               rm -f '/tmp/${GOOSE_BINARY}' &&
               sudo -n tar -xzf /tmp/journal-server-assets.tar.gz \
                 -C '${SERVER_INSTALL_DIR}' &&
               rm -f /tmp/journal-server-assets.tar.gz"

            ssh $SSH_OPTS "$DEPLOY_TARGET" \
              "SERVER_SERVICE='${SERVER_SERVICE}' \
               SERVER_INSTALL_DIR='${SERVER_INSTALL_DIR}' \
               SERVER_BINARY='${SERVER_BINARY}' \
               SERVER_ENV_FILE='${SERVER_ENV_FILE}' \
               DEPLOY_USER='${DEPLOY_USER}' sh -se" <<'REMOTE_SYSTEMD'
                SERVICE_UNIT="${SERVER_SERVICE%.service}.service"
                SERVICE_FILE="/etc/systemd/system/${SERVICE_UNIT}"

                sudo -n tee "${SERVICE_FILE}" >/dev/null <<SERVICE
[Unit]
Description=Journal API
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${DEPLOY_USER}
WorkingDirectory=${SERVER_INSTALL_DIR}
EnvironmentFile=${SERVER_ENV_FILE}
ExecStart=${SERVER_INSTALL_DIR}/${SERVER_BINARY}
Restart=on-failure
RestartSec=5
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
SERVICE

                sudo -n systemctl daemon-reload
                sudo -n systemctl enable "${SERVICE_UNIT}"
REMOTE_SYSTEMD

            if [ "$RUN_MIGRATIONS" = "true" ]; then
              ssh $SSH_OPTS "$DEPLOY_TARGET" \
                "SERVER_INSTALL_DIR='${SERVER_INSTALL_DIR}' \
                 SERVER_ENV_FILE='${SERVER_ENV_FILE}' \
                 GOOSE_INSTALL_PATH='${GOOSE_INSTALL_PATH}' sh -se" \
                <<'REMOTE_MIGRATIONS'
                  set -a
                  . "${SERVER_ENV_FILE}"
                  set +a
                  : "${DATABASE_URL:?DATABASE_URL is not set}"
                  "${GOOSE_INSTALL_PATH}" \
                    -dir "${SERVER_INSTALL_DIR}/migrations" \
                    postgres "${DATABASE_URL}" up
REMOTE_MIGRATIONS
            fi

            ssh $SSH_OPTS "$DEPLOY_TARGET" \
              "sudo -n systemctl restart '${SERVER_SERVICE}' &&
               sudo -n systemctl is-active --quiet '${SERVER_SERVICE}'"
          '''
        }
      }
    }

    stage('Deploy Client') {
      when {
        anyOf {
          branch 'main'
          expression { !env.BRANCH_NAME }
        }
      }
      steps {
        sshagent(credentials: [env.SSH_CREDENTIALS_ID]) {
          sh '''
            set -eu
            DEPLOY_TARGET="${DEPLOY_USER}@${DEPLOY_HOST}"
            SSH_OPTS="-o StrictHostKeyChecking=no"
            SSH_OPTS="$SSH_OPTS -o UserKnownHostsFile=/dev/null -o BatchMode=yes"

            tar -C "$CLIENT_DIR/dist" -czf "$ARTIFACT_DIR/client-dist.tar.gz" .
            scp $SSH_OPTS "$ARTIFACT_DIR/client-dist.tar.gz" \
              "$DEPLOY_TARGET:/tmp/journal-client-dist.tar.gz"

            ssh $SSH_OPTS "$DEPLOY_TARGET" \
              "sudo -n mkdir -p '${WEB_INSTALL_DIR}' &&
               sudo -n tar -xzf /tmp/journal-client-dist.tar.gz \
                 -C '${WEB_INSTALL_DIR}' &&
               rm -f /tmp/journal-client-dist.tar.gz &&
               sudo -n chown -R www-data:www-data '${WEB_INSTALL_DIR}' &&
               sudo -n systemctl reload nginx"
          '''
        }
      }
    }
  }

  post {
    always {
      archiveArtifacts(
        artifacts: '.local/jenkins-build/**',
        fingerprint: true,
        allowEmptyArchive: true
      )
    }
  }
}
