var test = new Vue({
    el: "#app",
    data: {
	on: false,
	observatory: observatoryBinding,
    },
    computed: {
	off: function(){ return !this.on;},
	toggleText: function(){return this.on ? "Stop" : "Start";},
    },
    methods: {
	toggleObserving: function(){
	    this.on = !this.on;
	},
	test: function() {
	    console.log(JSON.stringify(this.observatory));
	    console.log(this.observatory.data.TimeLine);
	},
    }
})
