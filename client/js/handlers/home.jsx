var React = require('react'),
	helpers = require('./helpers'),
	user = require('../user');

module.exports = function(ctx) {
	var avatarStyle = {
		background: 'url(' + user.get('avatar') + ')'
	};

	React.render(
		<section className='content'>
			Welcome <span style={avatarStyle} className='tinyAvatar' />{ user.get('username') }!
		</section>,
		helpers.getPageDiv()
	);
}
