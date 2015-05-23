var React = require('react'),
	user = require('../user'),
	Bot = require('./bot.jsx'),
	Greetings = require('./greetings.jsx'),
	Commands = require('./commands.jsx');

module.exports = React.createClass({
	render: function() {
		var avatarStyle = {
			background: 'url(' + user.get('avatar') + ')',
			'background-size': 'cover'
		};

		var channelLink = 'https://www.stream.me/' + user.get('username');

		return (
			<section className='content stall'>
				<p>Welcome 
					<a target='_blank' href={channelLink} title={user.get('username')}>
						<span style={avatarStyle} className='tinyAvatar' />
						{ user.get('username') }!
					</a>
				</p>
				<Bot />
				<Greetings />
				<Commands />
			</section>
		)
	}
});
