var React = require('react'),
	Bot = require('./bot.jsx'),
	Greetings = require('./greetings.jsx'),
	Commands = require('./commands.jsx');

module.exports = React.createClass({
	render: function() {
		return (
			<section className='content stall'>
				<Bot />
				<Greetings />
				<Commands />
			</section>
		)
	}
});
