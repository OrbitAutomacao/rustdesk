import ptBR from '@/utils/i18n/pt_BR.json'

export function T (key, params, num = 0) {
  const tran = ptBR[key]
  if (!tran) {
    return key
  }
  const msg = num > 1 ? (tran.Other ? tran.Other : tran.One) : tran.One
  //msg 是这样 {name} is name
  //params 是这样 {name: 'zhangsan'}
  //替换
  return msg.replace(/{(\w+)}/g, function (match, key) {
    return params?.[key] || match
  })
}
