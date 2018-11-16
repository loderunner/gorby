import {
  fetch
} from 'whatwg-fetch'

const mockRequest = {
  'id': 1,
  'timestamp': '2018-11-14T23:42:16.191251+01:00',
  'proto': 'HTTP/1.1',
  'method': 'POST',
  'host': 'google.com',
  'path': '/',
  'content_length': 9,
  'header': {
    'Accept': ['*/*'],
    'Accept-Encoding': ['gzip, deflate'],
    'Connection': ['keep-alive'],
    'Content-Length': ['9'],
    'Content-Type': ['application/x-www-form-urlencoded; charset=utf-8'],
    'User-Agent': ['HTTPie/0.9.9']
  },
  'body': 'dG90bz10b3Rv',
  'trailer': null,
  'query': {
    'hello': ['world']
  },
  'form': {
    'toto': ['toto']
  },
  'response': {
    'id': 1,
    'timestamp': '2018-11-14T23:42:16.299307+01:00',
    'proto': 'HTTP/1.1',
    'status': '405 Method Not Allowed',
    'status_code': 405,
    'content_length': 1601,
    'header': {
      'Allow': ['GET, HEAD'],
      'Content-Length': ['1601'],
      'Content-Type': ['text/html; charset=UTF-8'],
      'Date': ['Wed, 14 Nov 2018 22:42:16 GMT'],
      'Server': ['gws'],
      'X-Frame-Options': ['SAMEORIGIN'],
      'X-Xss-Protection': ['1; mode=block']
    },
    'body': 'PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ZW4+CiAgPG1ldGEgY2hhcnNldD11dGYtOD4KICA8bWV0YSBuYW1lPXZpZXdwb3J0IGNvbnRlbnQ9ImluaXRpYWwtc2NhbGU9MSwgbWluaW11bS1zY2FsZT0xLCB3aWR0aD1kZXZpY2Utd2lkdGgiPgogIDx0aXRsZT5FcnJvciA0MDUgKE1ldGhvZCBOb3QgQWxsb3dlZCkhITE8L3RpdGxlPgogIDxzdHlsZT4KICAgICp7bWFyZ2luOjA7cGFkZGluZzowfWh0bWwsY29kZXtmb250OjE1cHgvMjJweCBhcmlhbCxzYW5zLXNlcmlmfWh0bWx7YmFja2dyb3VuZDojZmZmO2NvbG9yOiMyMjI7cGFkZGluZzoxNXB4fWJvZHl7bWFyZ2luOjclIGF1dG8gMDttYXgtd2lkdGg6MzkwcHg7bWluLWhlaWdodDoxODBweDtwYWRkaW5nOjMwcHggMCAxNXB4fSogPiBib2R5e2JhY2tncm91bmQ6dXJsKC8vd3d3Lmdvb2dsZS5jb20vaW1hZ2VzL2Vycm9ycy9yb2JvdC5wbmcpIDEwMCUgNXB4IG5vLXJlcGVhdDtwYWRkaW5nLXJpZ2h0OjIwNXB4fXB7bWFyZ2luOjExcHggMCAyMnB4O292ZXJmbG93OmhpZGRlbn1pbnN7Y29sb3I6Izc3Nzt0ZXh0LWRlY29yYXRpb246bm9uZX1hIGltZ3tib3JkZXI6MH1AbWVkaWEgc2NyZWVuIGFuZCAobWF4LXdpZHRoOjc3MnB4KXtib2R5e2JhY2tncm91bmQ6bm9uZTttYXJnaW4tdG9wOjA7bWF4LXdpZHRoOm5vbmU7cGFkZGluZy1yaWdodDowfX0jbG9nb3tiYWNrZ3JvdW5kOnVybCgvL3d3dy5nb29nbGUuY29tL2ltYWdlcy9icmFuZGluZy9nb29nbGVsb2dvLzF4L2dvb2dsZWxvZ29fY29sb3JfMTUweDU0ZHAucG5nKSBuby1yZXBlYXQ7bWFyZ2luLWxlZnQ6LTVweH1AbWVkaWEgb25seSBzY3JlZW4gYW5kIChtaW4tcmVzb2x1dGlvbjoxOTJkcGkpeyNsb2dve2JhY2tncm91bmQ6dXJsKC8vd3d3Lmdvb2dsZS5jb20vaW1hZ2VzL2JyYW5kaW5nL2dvb2dsZWxvZ28vMngvZ29vZ2xlbG9nb19jb2xvcl8xNTB4NTRkcC5wbmcpIG5vLXJlcGVhdCAwJSAwJS8xMDAlIDEwMCU7LW1vei1ib3JkZXItaW1hZ2U6dXJsKC8vd3d3Lmdvb2dsZS5jb20vaW1hZ2VzL2JyYW5kaW5nL2dvb2dsZWxvZ28vMngvZ29vZ2xlbG9nb19jb2xvcl8xNTB4NTRkcC5wbmcpIDB9fUBtZWRpYSBvbmx5IHNjcmVlbiBhbmQgKC13ZWJraXQtbWluLWRldmljZS1waXhlbC1yYXRpbzoyKXsjbG9nb3tiYWNrZ3JvdW5kOnVybCgvL3d3dy5nb29nbGUuY29tL2ltYWdlcy9icmFuZGluZy9nb29nbGVsb2dvLzJ4L2dvb2dsZWxvZ29fY29sb3JfMTUweDU0ZHAucG5nKSBuby1yZXBlYXQ7LXdlYmtpdC1iYWNrZ3JvdW5kLXNpemU6MTAwJSAxMDAlfX0jbG9nb3tkaXNwbGF5OmlubGluZS1ibG9jaztoZWlnaHQ6NTRweDt3aWR0aDoxNTBweH0KICA8L3N0eWxlPgogIDxhIGhyZWY9Ly93d3cuZ29vZ2xlLmNvbS8+PHNwYW4gaWQ9bG9nbyBhcmlhLWxhYmVsPUdvb2dsZT48L3NwYW4+PC9hPgogIDxwPjxiPjQwNS48L2I+IDxpbnM+VGhhdOKAmXMgYW4gZXJyb3IuPC9pbnM+CiAgPHA+VGhlIHJlcXVlc3QgbWV0aG9kIDxjb2RlPlBPU1Q8L2NvZGU+IGlzIGluYXBwcm9wcmlhdGUgZm9yIHRoZSBVUkwgPGNvZGU+Lz9oZWxsbz13b3JsZDwvY29kZT4uICA8aW5zPlRoYXTigJlzIGFsbCB3ZSBrbm93LjwvaW5zPgo=',
    'trailer': null,
    'form': null
  }
}

const requests = [mockRequest]

const actions = {
  async listRequests({
    state,
    commit
  }, args) {
    const data = await fetch('http://localhost:8081/requests')
    commit()
  }
}

const getters = {
  requests: state => state
}

const mutations = {}

const constants = {
  // Actions
  ActionListRequests: 'listRequests',

  // Mutations
  MutationReceiveRequestsList: 'ReceiveRequestsList'
}

export default {
  state: requests,
  getters,
  actions,
  mutations,
  ...constants
}
