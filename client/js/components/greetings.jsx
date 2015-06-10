var React = require('react'),
	xhr = require('xhr'),
	mixins = require('./mixins');

module.exports = React.createClass({
	mixins: [mixins.errors, mixins.forms],

	getInitialState: function() {
		return {
			loading: false,
			newUser: "",
			returningUser: "",
			consecutiveUser: "",
			greetTrolls: false,
			answeringMachineOn: true,
			answeringMachine: ""
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
				greetTrolls: body.greetTrolls,
				answeringMachine: body.answeringMachine,
				answeringMachineOn: body.answeringMachineOn
			});
		}.bind(this));
	},

	save: function() {
		xhr({
			method: 'POST',
			uri: '/api/greeting-templates',
			json: {
				newUser: this.value('newUser'),
				returningUser: this.value('returningUser'),
				consecutiveUser: this.value('consecutiveUser'),
				answeringMachine: this.value('answeringMachine'),
				greetTrolls: this.state.greetTrolls,
				answeringMachineOn: this.state.answeringMachineOn
			}
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.saved();

			this.setState({
				loading: false,
				saved: true,
				newUser: body.newUser,
				returningUser: body.returningUser,
				answeringMachine: body.answeringMachine,
				consecutiveUser: body.consecutiveUser,
				greetTrolls: body.greetTrolls
			});
		}.bind(this));
	},

	componentWillMount: function() {
		this.get();
	},

	toggleGreetTrolls: function() {
		this.setState({
			greetTrolls: !this.state.greetTrolls
		}, this.save);
	},

	toggleAnsweringMachine: function() {
		this.setState({
			answeringMachineOn: !this.state.answeringMachineOn
		}, this.save);
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

		var trollClass, trollButtonText, answeringMachineOnOff;
		if(this.state.greetTrolls) {
			trollClass = 'radio-button on';
			trollButtonText = 'Stop !meeping with anonymous users';
		} else {
			trollClass = 'radio-button off';
			trollButtonText = 'Start !meeping with anonymous users';
		}

		if(this.state.answeringMachineOn) {
			answeringMachineOnOff = <a onClick={this.toggleAnsweringMachine} className='on fr'>Turn answering machine off</a>
		} else {
			answeringMachineOnOff = <a onClick={this.toggleAnsweringMachine} className='off fr'>Turn answering machine on</a>
		}

		return (
			<form>
				<div className='field'>
					<label>New User Greeting</label>
					<textarea maxLength='250' name='newUser' placeholder='New user greeting' defaultValue={this.state.newUser}></textarea>
				</div>
				<div className='field'>
					<label>Returning User Greeting</label>
					<textarea maxLength='250' name='returningUser' placeholder='Returning user greeting' defaultValue={this.state.returningUser}></textarea>
				</div>
				<div className='field'>
					<label>Consecutive User Greeting</label>
					<textarea maxLength='250' name='consecutiveUser' placeholder='Consecutive user greeting' defaultValue={this.state.consecutiveUser}></textarea>
				</div>
				<div className='field'>
					<label>Answering Machine {answeringMachineOnOff}</label>
					<textarea maxLength='250' name='answeringMachine' placeholder='Answering machine greeting' defaultValue={this.state.answeringMachine}></textarea>
				</div>
				<div className='field'>
					<label>
						<a onClick={this.toggleGreetTrolls} className={trollClass} title={trollButtonText}>
							{trollButtonText}
						</a>
					</label>
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
