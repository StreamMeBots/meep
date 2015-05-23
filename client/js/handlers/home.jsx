var React = require('react'),
	helpers = require('./helpers'),
	user = require('../user'),
	Bot = require('../components/bot.jsx'),
	Greetings = require('../components/greetings.jsx');

module.exports = function(ctx) {
	var avatarStyle = {
		background: 'url(' + user.get('avatar') + ')',
		'background-size': 'cover'
	};

	var channelLink = 'https://www.stream.me/' + user.get('username');

	React.render(
		<section className='content stall'>
			<p>Welcome 
				<a target='_blank' href={channelLink} title={user.get('username')}>
					<span style={avatarStyle} className='tinyAvatar' />
					{ user.get('username') }!
				</a>
			</p>
			<Bot />
			<Greetings />
		</section>,
		helpers.getPageDiv()
	);
}
