import request from '@/utils/request'

export function list (params) {
  return request({
    url: '/my/peer/list',
    params,
  })
}

export function claim (rowId) {
  return request({
    url: '/my/peer/claim',
    method: 'post',
    data: { row_id: rowId },
  })
}
