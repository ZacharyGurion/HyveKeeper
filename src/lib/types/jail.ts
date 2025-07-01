export type JailStatus = 'running' | 'stopped' | 'paused' | 'error';

export interface Jail {
    id: number;
    name: string;
    status: JailStatus;
    memory: number;
    cpu: number;
    disk: number;
    uptime: string;
    ip: string;
}
