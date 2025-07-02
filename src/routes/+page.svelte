<script lang="ts">
  import type { VM, NodeStats } from '$lib/types/vm';
	import type { Jail } from '$lib/types/jail';
  import { onMount } from 'svelte';
  let nodeStats: NodeStats = {
    cpu: 35,
    memory: { used: 32768, total: 65536 },
    storage: { used: 500, total: 2000 },
    uptime: '45d 12h'
  };
  let vms: VM[] = [];
	let jails: Jail[] = [];

  onMount(async () => {
    try {
      const response = await fetch('https://api.zgurion.com/vms');
      vms = await response.json();
    } catch (error) {
      console.error('Error fetching VMs:', error);
    }
    try {
      const response = await fetch('https://api.zgurion.com/jails');
			const jailData = await response.json();
      jails = jailData.jails;
			console.log(jails);
    } catch (error) {
      console.error('Error fetching Jails:', error);
    }
  });



  function getStatusColor(status: string): string {
    return status === 'running' ? 'var(--green)' : 'var(--red)';
  }

  function mByteFmt(size: number): string {
    if (size < 1024) return `${size.toFixed(2)} MB`;
    return `${(size / 1024).toFixed(2)} GB`;
  }

  function byteFmt(size: number): string {
    if (size < 1024*1000) return `${(size/1024).toFixed(2)} KB`;
    if (size < 1024*1000000) return `${(size/(1024*1024)).toFixed(2)} MB`;
    return `${(size / (1024*1024*1024)).toFixed(2)} GB`;
  }

  function formatPercentage(used: number, total: number): number {
    return Math.round((used / total) * 100);
  }
</script>

<div class="py-6">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="p-6 mb-6" style="background-color: var(--background-1); box-shadow: 0 4px 6px var(--shadow-color); border-radius: var(--radius-lg);">
      <h2 class="text-2xl font-semibold mb-4" style="color: var(--foreground-0);">Node Status</h2>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <div class="text-sm font-medium" style="color: var(--foreground-1);">CPU Usage</div>
          <div class="mt-1">
            <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
              <span>{nodeStats.cpu}%</span>
            </div>
            <div class="w-full h-2.5" style="background-color: var(--background-2); border-radius: var(--radius-full);">
              <div class="h-2.5"
                 style="width: {nodeStats.cpu}%; background-color: var(--blue); border-radius: var(--radius-full);">
              </div>
            </div>
          </div>
        </div>

        <div>
          <div class="text-sm font-medium" style="color: var(--foreground-1);">Memory</div>
          <div class="mt-1">
            <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
              <span>{mByteFmt(nodeStats.memory.used)} / {mByteFmt(nodeStats.memory.total)}</span>
            </div>
            <div class="w-full h-2.5" style="background-color: var(--background-2); border-radius: var(--radius-full);">
              <div class="h-2.5"
                 style="width: {formatPercentage(nodeStats.memory.used, nodeStats.memory.total)}%; background-color: var(--purple); border-radius: var(--radius-full);">
              </div>
            </div>
          </div>
        </div>

        <div>
          <div class="text-sm font-medium" style="color: var(--foreground-1);">Storage</div>
          <div class="mt-1">
            <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
              <span>{mByteFmt(nodeStats.storage.used)} / {mByteFmt(nodeStats.storage.total)}</span>
            </div>
            <div class="w-full h-2.5" style="background-color: var(--background-2); border-radius: var(--radius-full);">
              <div class="h-2.5"
                 style="width: {formatPercentage(nodeStats.storage.used, nodeStats.storage.total)}%; background-color: var(--orange); border-radius: var(--radius-full);">
              </div>
            </div>
          </div>
        </div>

        <div>
          <div class="text-sm font-medium" style="color: var(--foreground-1);">Uptime</div>
          <div class="text-lg font-semibold" style="color: var(--foreground-0);">{nodeStats.uptime}</div>
        </div>
      </div>
    </div>
    <!--  Jails  Section -->
    <div class="p-6" style="background-color: var(--background-1); box-shadow: 0 4px 6px var(--shadow-color); border-radius: var(--radius-lg);">
      <h2 class="text-2xl font-semibold mb-4" style="color: var(--foreground-0);">Jails</h2>
      <div class="grid gap-4">
        {#each jails as jail (jail.id)}
          <div class="p-4" style="background-color: var(--background-2); border-radius: var(--radius-md);">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-lg font-semibold" style="color: var(--foreground-0);">{jail.name}</h3>
                <div class="flex items-center mt-1">
                  <div class="w-2 h-2 rounded-full mr-2" style="background-color: {getStatusColor(jail.status)}"></div>
                  <span class="text-sm" style="color: var(--foreground-1);">{jail.status}</span>
                  <span class="text-sm mx-2" style="color: var(--foreground-2);">•</span>
                  <span class="text-sm" style="color: var(--foreground-1);">{jail.ip}</span>
                </div>
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <!-- CPU -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">
                  CPU
                </div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{jail.cpu}%</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {jail.cpu}%; background-color: var(--blue); border-radius: var(--radius-full);"></div>
                </div>
              </div>

              <!-- Memory -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">Memory</div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{byteFmt(jail.memory)} / {byteFmt(190000000)}</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {formatPercentage(jail.memory, jail.memory)}%; background-color: var(--purple); border-radius: var(--radius-full);"></div>
                </div>
              </div>

              <!-- Disk -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">Disk</div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{jail.disk}GB / {jail.disk}GB</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {formatPercentage(jail.disk, jail.disk)}%; background-color: var(--orange); border-radius: var(--radius-full);"></div>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
    <!-- Virtual Machines Section -->
    <div class="p-6" style="background-color: var(--background-1); box-shadow: 0 4px 6px var(--shadow-color); border-radius: var(--radius-lg);">
      <h2 class="text-2xl font-semibold mb-4" style="color: var(--foreground-0);">Virtual Machines</h2>
      <div class="grid gap-4">
        {#each vms as vm (vm.id)}
          <div class="p-4" style="background-color: var(--background-2); border-radius: var(--radius-md);">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-lg font-semibold" style="color: var(--foreground-0);">{vm.name}</h3>
                <div class="flex items-center mt-1">
                  <div class="w-2 h-2 rounded-full mr-2" style="background-color: {getStatusColor(vm.status)}"></div>
                  <span class="text-sm" style="color: var(--foreground-1);">{vm.status}</span>
                  <span class="text-sm mx-2" style="color: var(--foreground-2);">•</span>
                  <span class="text-sm" style="color: var(--foreground-1);">{vm.ip}</span>
                </div>
              </div>
              <div class="text-sm" style="color: var(--foreground-1);">
                Uptime: {vm.uptime}
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <!-- CPU -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">
                  CPU ({vm.cpu.cores} cores)
                </div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{vm.cpu.usage}%</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {vm.cpu.usage}%; background-color: var(--blue); border-radius: var(--radius-full);"></div>
                </div>
              </div>

              <!-- Memory -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">Memory</div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{mByteFmt(vm.memory.used)} / {mByteFmt(vm.memory.total)}</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {formatPercentage(vm.memory.used, vm.memory.total)}%; background-color: var(--purple); border-radius: var(--radius-full);"></div>
                </div>
              </div>

              <!-- Disk -->
              <div>
                <div class="text-sm font-medium mb-1" style="color: var(--foreground-1);">Disk</div>
                <div class="flex justify-between text-sm mb-1" style="color: var(--foreground-2);">
                  <span>{vm.disk.used}GB / {vm.disk.total}GB</span>
                </div>
                <div class="w-full h-2.5" style="background-color: var(--background-1); border-radius: var(--radius-full);">
                  <div class="h-2.5" style="width: {formatPercentage(vm.disk.used, vm.disk.total)}%; background-color: var(--orange); border-radius: var(--radius-full);"></div>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>
</div>
