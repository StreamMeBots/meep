var React = require('react'),
	moment = require('moment'),
	xhr = require('xhr'),
	mixins = require('./mixins'),
	OneMinute = 1000 * 60;

module.exports = React.createClass({
	mixins: [mixins.errors],
	getInitialState: function() {
		return {
			starting: false,
			loading: true,
			state: null,
			started: null
		}
	},

	// Get fresh data
	get: function() {
		this.getting = true;

		xhr({
			method: 'GET',
			uri: '/api/bot',
			json: true
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.getting = false;

			this.setState({
				loading: false,
				starting: false,
				state: body.State,
				started: moment(body.Started)
			});
		}.bind(this));
	},

	shouldGet: function() {
		return !this.getting;
	},

	// Try to get fresh data
	load: function() {
		if(!this.shouldGet()) {
			return;
		}
		this.get();
	},

	// Start the meeping!
	start: function() {
		this.setState({
			starting: true
		});

		xhr({
			method: 'POST',
			uri: '/api/bot',
			json: true
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.get();
		}.bind(this));
	},

	// Stop the meeping!
	stop: function() {
		this.setState({
			stopping: true
		});

		xhr({
			method: 'DELETE',
			uri: '/api/bot',
			json: true
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.setState({
				loading: false
			});
			this.get();
		}.bind(this));
	},

	// Start an infrequent short poll
	componentWillMount: function() {
		this.get();
		this.slowShortPoll = setInterval(this.get.bind(this), OneMinute);
	},

	// Stop the short poll
	componentWillUnmount: function() {
		clearInterval(this.slowShortPoll);
	},

	// Contents
	getContents: function() {
		if(this.state.error) {
			return (<p className='error'>{this.state.error}</p>);
		}

		if (this.state.loading) {
			return (
				<div>
					<div className='details'>Trying to load the !meeping details</div>
					<div className='actions'>
						<a className='button' onClick={this.load} title='Refresh'>Check !meep</a>
					</div>
				</div>
			);
		}

		if (this.state.starting) {
			return (
				<div>
					<div className='details'>Starting !meep</div>
					<div className='actions'>
						<a className='button' onClick={this.load} title='Refresh'>Check !meep</a>
					</div>
				</div>
			);
		}

		switch(this.state.state) {
			case 'notStarted': {
				return (
					<div>
						<div className='details'>Not started</div>
						<div className='actions'>
							<a className='button' onClick={this.start} title='Start'>Start !meep</a>
						</div>
					</div>
				)
			}
			case 'Joined': {
				return (
					<div>
						<div className='details'>!meeping since {this.state.started.format('dddd M/DD h:mm a')}</div>
						<div className='actions'>
							<a className='button' onClick={this.stop} title='Stop'>Stop !meeping</a>
						</div>
					</div>
				)
			}
			case 'Connecting': {
				return (
					<div>
						<div className='details'>!meep is connecting to your channel</div>
						<div className='actions'>
							<a className='button' onClick={this.stop} title='Stop'>Stop !meeping</a>
						</div>
					</div>
				)
			}
		}
	},

	// Main render
	render: function() {
		var contents = this.getContents();

		return (
			<section className='card'>
				<h1>!meeping status</h1>
				{contents}
			</section>
		);
	}
})