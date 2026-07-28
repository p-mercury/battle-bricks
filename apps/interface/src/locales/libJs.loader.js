import { registerLoaders } from 'wuchale/load-utils'
import { loadCatalog, loadCount } from './.wuchale/libJs.proxy.js'

const key = 'libJs'

// two exports. can be used anywhere
export const getRuntime = registerLoaders(key, loadCatalog, loadCount)
export const getRuntimeRx = getRuntime
