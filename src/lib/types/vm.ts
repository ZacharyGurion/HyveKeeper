export type VMStatus = 'running' | 'stopped' | 'paused' | 'error';

export interface VM {
    id: number;
    name: string;
    status: VMStatus;
    memory: {
        used: number;
        total: number;
    };
    cpu: {
        cores: number;
        usage: number;
    };
    disk: {
        used: number;
        total: number;
    };
    uptime: string;
    ip: string;
}

export interface NodeStats {
    cpu: number;
    memory: {
        used: number;
        total: number;
    };
    storage: {
        used: number;
        total: number;
    };
    uptime: string;
}