<template>
  <b-table :data="requests" :selected.sync="selected" narrowed hoverable>
    <template slot-scope="props">
      <b-table-column field="id" label="ID" numeric>{{ props.row.id }}</b-table-column>
      <b-table-column
        field="timestamp"
        label="Time"
      >{{ props.row.timestamp | moment('HH:mm:ss.SSS') }}</b-table-column>
      <b-table-column field="method" label="Method">{{ props.row.method }}</b-table-column>
      <b-table-column field="host" label="Host">{{ props.row.host }}</b-table-column>
      <b-table-column field="path" label="Path" width="6000">{{ props.row.path }}</b-table-column>
    </template>
  </b-table>
</template>

<script>
import { Table } from 'buefy'
import { mapGetters } from 'vuex'
import Requests from '../store/modules/requests'

export default {
  name: 'RequestsTable',
  components: {
    Table
  },
  computed: {
    ...mapGetters(['requests'])
  },
  data() {
    return {
      selected: undefined
    }
  },
  beforeCreate() {
    this.storeSubscribe = this.$store.subscribe((mutation, state) => {
      switch (mutation) {
        case Requests.MutationReceiveRequestsList:
        break
      }
    })
    this.$store.dispatch(Requests.ActionSubscribe)
  }
}
</script>
