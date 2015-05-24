var mixins = module.exports = {};

mixins.errors = {
	// Handle errors, set a message state and log them
	error: function(err) {
		console.log('bot:error', err);
		this.setState({
			error: 'Sorry, there was an internal server error.  Please try again later.'
		});
	}
}

mixins.forms = {
	value: function(v) {
		return this.getDOMNode().querySelector('[name="'+v+'"]').value;
	},

	saved: function() {
		if(this.savedTO) {
			clearTimeout(this.savedTO);
		}
		this.savedTO = setTimeout(function(){
			this.setState({
				saved: false
			});
		}.bind(this), 5000);
	}
};
