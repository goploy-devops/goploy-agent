import { Request } from './types'

export class General extends Request {
  readonly url = '/general'
  readonly method = 'get'
  public datagram!: {
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
  public datagram!: {
    avg: string
    avg5: string
    avg15: string
    cores: string
  }
}

export class RAM extends Request {
  readonly url = '/ram'
  readonly method = 'get'
  public datagram!: {
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
