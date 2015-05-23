module.exports = {
	// Handle errors, set a message state and log them
	error: function(err) {
		console.log('bot:error', err);
		this.setState({
			error: 'Sorry, there was an internal server error.  Please try again later.'
		});
	}
}
