import { fetch } from 'whatwg-fetch'
import moment from 'vue-moment'

// Constants
const ActionListRequests = 'listRequests'
const ActionSubscribe = 'subscribe'

const MutationReceiveRequestsList = 'ReceiveRequestsList'
const MutationReceiveRequest = 'ReceiveRequest'

const constants = {
  // Actions
  ActionListRequests,
  ActionSubscribe,

  // Mutations
  MutationReceiveRequestsList,
  MutationReceiveRequest
}

// Actions
const actions = {
  async listRequests({ commit }) {
    try {
      const res = await fetch('http://localhost:8081/requests')
      if (res.ok) {
        const data = await res.json()
        commit(MutationReceiveRequestsList, { requests: data })
      } else {
        const errMsg = res.text()
        console.error(errMsg)
      }
    } catch (err) {
      console.error(err)
    }
  },
  subscribe({ commit }, args) {
    if (listener) {
      return
    }
    try {
      let url = 'http://localhost:8081/requests'
      if (args && args.start) {
        const start = moment(args.start)
        url += `?start=${start.toISOString()}`
      }
      listener = new EventSource(url)
      listener.onmessage = (evt) => {
        try {
          const data = JSON.parse(evt.data)
          commit(MutationReceiveRequest, { request: data.request })
        } catch (err) {
          console.error(err)
        }
      }
    } catch (err) {
      console.error(err)
      listener = null
    }
  }
}

// State
const requests = []

let listener = null

// Getters
const getters = {
  requests: state => state
}

// Mutations
const mutations = {
  [MutationReceiveRequestsList]: (state, { requests }) => {
    state = requests.map(({ request, response }) => ({ ...request, response }))
  },
  [MutationReceiveRequest]: (state, { request }) => {
    state.push(request)
  }
}

export default {
  state: requests,
  getters,
  actions,
  mutations,
  ...constants
}
