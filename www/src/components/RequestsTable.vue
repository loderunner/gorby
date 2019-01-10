<template>
  <b-table :data="requests" :selected.sync="selected" @click="rowSelected" narrowed hoverable>
    <template slot-scope="props">
      <!-- <b-table-column field="id" label="ID" numeric>{{ props.row.id }}</b-table-column> -->
      <b-table-column
        field="timestamp"
        label="Time"
      >{{ props.row.timestamp | moment('HH:mm:ss.SSS') }}</b-table-column>
      <b-table-column field="host" label="Host">{{ props.row.host }}</b-table-column>
      <b-table-column field="path" label="Path">{{ props.row.path }}</b-table-column>
      <b-table-column field="status" label="Status">
        <LoadingAnimation v-if="!(props.row.response)"/>
        <b-tooltip :label="props.row.response.status" v-else>{{ props.row.response.status_code }}</b-tooltip>
      </b-table-column>
      <b-table-column field="size" label="Size">
        <span
          v-if="!!(props.row.response)"
        >{{ props.row.response.content_length | numeralFormat('0.[0] b') }}</span>
      </b-table-column>
    </template>
  </b-table>
</template>

<script>
import { Table } from 'buefy'
import { mapGetters } from 'vuex'
import LoadingAnimation from './LoadingAnimation'
import Requests from '../store/modules/requests'

export default {
  name: 'RequestsTable',
  components: {
    Table,
    LoadingAnimation
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
  },
  methods: {
    rowSelected(row) {
      if (row === this.selected) {
        this.selected = undefined
      }
    }
  }
}
</script>
