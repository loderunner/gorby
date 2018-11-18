import { fetch } from 'whatwg-fetch'

// Constants
const ActionListRequests = 'listRequests'
const MutationReceiveRequestsList = 'ReceiveRequestsList'

const constants = {
  // Actions
  ActionListRequests,

  // Mutations
  MutationReceiveRequestsList
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
  }
}

// State
const requests = []

    // Getters
    const getters = {
      requests: state => state
    }

// Mutations
const mutations = {
  [MutationReceiveRequestsList]: (state, { requests }) => {
    state.requests =
        requests.map(({ request, response }) => ({ ...request, response }))
  }
}

export default {
  state: { requests }, getters, actions, mutations, ...constants
}
