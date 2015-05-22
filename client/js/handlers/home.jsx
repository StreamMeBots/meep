var React = require('react'),
	helpers = require('./helpers'),
	user = require('../user'),
	Bot = require('../components/bot.jsx');

module.exports = function(ctx) {
	var avatarStyle = {
		background: 'url(' + user.get('avatar') + ')'
	};

	React.render(
		<section className='content'>
			<p>Welcome <span style={avatarStyle} className='tinyAvatar' />{ user.get('username') }!</p>
			<Bot />
		</section>,
		helpers.getPageDiv()
	);
}
