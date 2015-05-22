var React = require('react'),
	moment = require('moment'),
	xhr = require('xhr'),
	OneMinute = 1000 * 60;

module.exports = React.createClass({
	getInitialState: function() {
		return {
			starting: false,
			loaded: false,
			state: null,
			started: null
		}
	},

	error: function(err) {
		console.error(err);
		this.setState({
			error: 'Sorry, we could not load your bot information.  Please try again later.'
		});
	},

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
				loaded: true,
				starting: false,
				state: body.State,
				started: moment(body.Started)
			});
		}.bind(this));
	},

	shouldGet: function() {
		return !this.getting;
	},

	load: function() {
		if(!this.shouldGet()) {
			return;
		}
		this.get();
	},

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
				loaded: true
			});
			this.get();
		}.bind(this));
	},

	componentWillMount: function() {
		this.get();
		this.slowShortPoll = setInterval(this.get.bind(this), OneMinute);
	},

	componentWillUnmount: function() {
		clearInterval(this.slowShortPoll);
	},

	render: function() {
		var contents;

		if(this.state.error) {
			contents = (<p className='error'>{this.state.error}</p>);
		} else if (this.state.starting) {
			contents = (
				<div>
					<div className='details'>Starting !meep</div>
					<div className='actions'>
						<a className='button' onClick={this.load} title='Refresh'>Check !meep</a>
					</div>
				</div>
			);
		} else {
			switch(this.state.state) {
				case 'notStarted': {
					contents = (
						<div>
							<div className='details'>Not started</div>
							<div className='actions'>
								<a className='button' onClick={this.start} title='Start'>Start !meep</a>
							</div>
						</div>
					)
					break;
				}
				case 'Joined': {
					contents = (
						<div>
							<div className='details'>!meep is in your channel !meeping</div>
							<div className='actions'>
								<a className='button' onClick={this.stop} title='Stop'>Stop !meeping</a>
							</div>
						</div>
					)
					break;
				}
			}
		}


		return (
			<section className='card'>
				<h1>!meep bot status</h1>
				{contents}
			</section>
		);
	}
})