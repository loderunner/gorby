<template>
  <div id="requestsView" class="container">
    <RequestsTable/>
  </div>
</template>
<script>
import RequestsTable from '@/components/RequestsTable'
export default {
    name: 'Requests',
    components: {
        RequestsTable
    }
}
</script>
<style lang="scss" scoped>
#requestsView {
  max-width: 800px;
}
</style>

