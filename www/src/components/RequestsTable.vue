<template>
  <b-table :data="data" :columns="columns"></b-table>
</template>

<script>
import { Table } from 'buefy'

export default {
  name: 'RequestsTable',
  components: {
    Table
  },
  data() {
    return {
      data: [
        {time: new Date(), name: 'AAAAAA'},
        {time: new Date(), name: 'AAAAAB'},
      ],
      columns: [
        {
          field: 'time',
          label: 'Time',
        },
        {
          field: 'name',
          label: 'Name',
        }
      ]
    }
  }
}
</script>


