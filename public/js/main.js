setInterval(() => {
  dash.innerHTML = dash.innerHTML === "" ? "_" : "";
}, 500);

app = new Vue({
  el: "#app",
  data: {
    serverTime: "flight router monitor",
    status: "disconnected",
    numFlights: 0,
    cacheHits: 0,
    cacheMisses: 0,
  },
  mounted: function () {
    this.init();
  },
  methods: {
    init() {
      var socket = new WebSocket("ws://localhost:8080/ws/status");
      socket.onopen = this.connected;
      socket.onclose = this.disconnected;
      socket.onmessage = this.update;
    },
    update(e) {
      data = JSON.parse(e.data);
      this.numFlights = data.numberFlights;
      this.cacheHits = data.cache.hits;
      this.cacheMisses = data.cache.misses;
      this.serverTime = data.time;
    },
    connected() {
      this.status = "up";
    },
    disconnected() {
      this.status = "disconnected";
    },
  },
});
