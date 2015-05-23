var React = require('react'),
	xhr = require('xhr'),
	errors = require('./errors');

module.exports = React.createClass({
	mixins: [errors],

	getInitialState: function() {
		return {
			loading: false,
			newUser: "",
			returningUser: "",
			consecutiveUser: "",
			greetTrolls: false
		}
	},

	get: function() {
		if(this.state.loading) {
			return;
		}

		this.setState({
			loading: true
		});

		xhr({
			method: 'GET',
			url: '/api/greeting-templates',
			json: true
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.setState({
				loading: false,
				newUser: body.newUser,
				returningUser: body.returningUser,
				consecutiveUser: body.consecutiveUser,
				greetTrolls: body.greetTrolls
			});
		}.bind(this));
	},

	save: function() {
		var value = function(v) {
			return this.getDOMNode().querySelector('[name="'+v+'"]').value;
		}.bind(this);

		xhr({
			method: 'POST',
			uri: '/api/greeting-templates',
			json: {
				newUser: value('newUser'),
				returningUser: value('returningUser'),
				consecutiveUser: value('consecutiveUser'),
			}
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			if(this.savedTO) {
				clearTimeout(this.savedTO);
			}
			this.savedTO = setTimeout(function(){
				this.setState({
					saved: false
				});
			}.bind(this), 5000);

			this.setState({
				loading: false,
				saved: true,
				newUser: body.newUser,
				returningUser: body.returningUser,
				consecutiveUser: body.consecutiveUser,
				greetTrolls: body.greetTrolls
			});
		}.bind(this));
	},

	componentWillMount: function() {
		this.get();
	},

	getContents: function() {
		if(this.state.error) {
			return (<p className='error'>{this.state.error}</p>);
		}

		if(this.state.loading) {
			return (
				<div>
					<div className='details'>Trying to load the !meeping details</div>
					<div className='actions'>
						<a className='button' onClick={this.get} title='Refresh'>Check !meep</a>
					</div>
				</div>
			);
		}

		var savedString = this.state.saved ? (<span className='message'>Saved!</span>) : '';

		return (
			<form>
				<div className='field'>
					<label>New User Greeting</label>
					<textarea maxLength='250' name='newUser' placeholder='New user greeting'>{this.state.newUser}</textarea>
				</div>
				<div className='field'>
					<label>Returning User Greeting</label>
					<textarea maxLength='250' name='returningUser' placeholder='Returning user greeting'>{this.state.returningUser}</textarea>
				</div>
				<div className='field'>
					<label>Consecutive User Greeting</label>
					<textarea maxLength='250' name='consecutiveUser' placeholder='Consecutive user greeting'>{this.state.consecutiveUser}</textarea>
				</div>
				<div className='actions'>
					{savedString}
					<a className='button' onClick={this.save} title='Save'>Save the !meeping greetings</a>
				</div>
			</form>
		)
	},

	render: function() {
		var contents = this.getContents();

		return (
			<section className='card'>
				<h1>!meeping greetings</h1>
				{contents}
			</section>
		);
	}
})