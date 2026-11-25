import { afterEach, expect, it, vi } from 'vitest'
import { listParentMissedTopics } from './parentMissedTopics'
afterEach(()=>vi.restoreAllMocks())
it('requests the selected child and academic year',async()=>{
 const fetch=vi.spyOn(globalThis,'fetch').mockResolvedValue(new Response(JSON.stringify({items:[]})))
 expect(await listParentMissedTopics('child-1','year-1')).toEqual([])
 expect(fetch).toHaveBeenCalledWith('/api/v1/parent/students/child-1/missed-topics?academicYearId=year-1',expect.objectContaining({credentials:'include'}))
})
