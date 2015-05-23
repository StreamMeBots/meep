var React = require('react'),
	xhr = require('xhr'),
	mixins = require('./mixins'),
	pubsub = require('../pubsub');


var Command = React.createClass({
	mixins: [mixins.errors, mixins.forms],

	getInitialState: function() {
		return {
			name: this.props.name,
			template: this.props.template,
			deleted: false
		};
	},

	del: function(name, prompted) {
		if(!name || (this.lastName === this.newName && !prompted)) {
			return;
		}

		xhr({
			method: 'DELETE',
			uri: '/api/commands/' + name,
			json: true
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			if(prompted) {
				this.setState({
					deleted: true
				});
			}

			if(prompted) {
				pubsub.emit('command:updated', {
					id: this.props.id,
					deleted: true
				});
			}
		}.bind(this));
	},

	save: function() {
		if(!this.value('name')) {
			return;
		}

		if(this.props.name) {
			this.del(this.props.name);
		}

		var d = {
			name: this.value('name'),
			template: this.value('template')
		};

		this.lastName = this.newName || this.props.name;
		this.newName = d.name;

		xhr({
			method: 'PUT',
			uri: '/api/commands',
			json: d
		}, function(err, resp, body) {
			if(err) {
				return this.error(err);
			}
			if(resp.statusCode !== 200) {
				return this.error(new Error('unexpected status code'));
			}

			this.saved();

			pubsub.emit('command:updated', {
				id: this.props.id,
				name: d.name,
				template: d.template
			});

			this.setState({
				loading: false,
				saved: true
			});
		}.bind(this));
	},

	clickDel: function(e) {
		var answer = confirm('meeping delete it?!');

		if(!answer) {
			return;
		}

		this.del(this.value('name'), true);
	},

	render: function() {
		if(this.state.deleted) {
			return <p>Deleted!</p>;
		}
		var savedString = this.state.saved ? (<span className='message'>Saved!</span>) : '';

		return (
			<form>
				<div className='field'>
					<label>Command</label>
					<input maxLength='250' name='name' placeholder='Command' type='text' defaultValue={this.state.name} />
				</div>

				<div className='field'>
					<label>Response</label>
					<textarea maxLength='250' name='template' placeholder='Response'>{this.state.template}</textarea>
				</div>

				<div className='actions'>
					{savedString}
					<a className='button chill' onClick={this.clickDel} title='Delete'>Delete</a>
					<a className='button' onClick={this.save} title='Save'>Save the !meeping greetings</a>
				</div>
			</form>
		)
	}
});

module.exports = React.createClass({
	mixins: [mixins.errors],

	getInitialState: function() {
		return {
			loading: false,
			commands: []
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
			url: '/api/commands',
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
				commands: body
			});
		}.bind(this));
	},

	updateCommand: function(command) {
		var pushOne=true, c;

		for(var i=0; i<this.state.commands.length; i++) {
			c = this.state.commands[i];
			if(command.id === c.id) {
				c.name = command.deleted ? 'deleted' : command.name;
				c.template = command.deleted ? 'deleted' : command.template;
			}

			if(!c.name && !c.template) {
				pushOne = false;
			}
		}

		if(!pushOne) {
			return;
		}

		this.state.commands.push({
			name: '',
			template: ''
		});
		this.setState({
			commands: this.state.commands
		});
	},

	componentWillMount: function() {
		this.get();

		pubsub.on('command:updated', this.updateCommand.bind(this));
	},

	componentWillUnmount: function() {
		pubsub.off('command:updated', this.updateCommand.bind(this));
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

		var commands = [], c, pushOne = true;
		var add = function(command) {
			if(commands.length > 0) {
				commands.push( <div className='separator' /> );
			}
			commands.push( <Command id={command.id} name={command.name} template={command.template} /> );
		}
		for(var i=0; i<this.state.commands.length; i++) {
			c = this.state.commands[i];

			if(!c.name && !c.template) {
				pushOne = false;
			}
			c.id = i;
			add(c);
		}

		if(pushOne) {
			var d = {id: i, name: '', template: ''};
			this.state.commands.push(d)
			add(d);
		}

		return {commands};
	},

	render: function() {
		var contents = this.getContents();

		return (
			<section className='card'>
				<h1>!meeping commands</h1>
				{contents}
			</section>
		);
	}
})