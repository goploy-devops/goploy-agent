import { Request, Pagination } from './types'

export class Chart extends Request {
  readonly url = '/chart'
  readonly method = 'get'
  public param: {
    type: number
    datetimeRange: string
  }
  constructor(param: Chart['param']) {
    super()
    this.param = param
  }
}


export class General extends Request {
  readonly url = '/general'
  readonly method = 'get'
  public declare datagram: {
    kernelVersion: string
    os: string
    cores: string
    hostname: string
    uptime: string
  }
}

export class Loadavg extends Request {
  readonly url = '/loadavg'
  readonly method = 'get'
  public declare datagram: {
    avg: string
    avg5: string
    avg15: string
    cores: string
  }
}

export class RAM extends Request {
  readonly url = '/ram'
  readonly method = 'get'
  public declare datagram: {
    total: number
    free: number
  }
}

export class CPU extends Request {
  readonly url = '/cpu'
  readonly method = 'get'
}

export class Net extends Request {
  readonly url = '/net'
  readonly method = 'get'
}

export class DiskUsage extends Request {
  readonly url = '/diskUsage'
  readonly method = 'get'
}

export class DiskIOStat extends Request {
  readonly url = '/diskIOStat'
  readonly method = 'get'
}

export class CronData {
  public datagram!: {
    id: number
    serverId: number
    expression: string
    command: string
    singleMode: number
    logLevel: number
    description: string
    creator: string
    creatorId: number
    editor: string
    editorId: number
    state: number
    InsertTime: string
    UpdateTime: string
  }
}

export class CronList extends Request {
  readonly url = '/cronList'
  readonly method = 'get'
  public declare datagram: CronData['datagram'][]
}

export class CronLogData {
  public datagram!: {
    id: number
    message: string
    execCode: number
    reportTime: string
  }
}

export class CronLogs extends Request {
  readonly url = '/cronLogs'
  readonly method = 'get'
  public pagination: Pagination
  public param: {
    id: number
  }

  public declare datagram: CronLogData['datagram'][]

  constructor(param: CronLogs['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}
