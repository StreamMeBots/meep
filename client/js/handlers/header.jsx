var React = require('react'),
	helpers = require('./helpers'),
	user = require('../user');

module.exports = function(el) {
	var menu = '';

	if(user.isAuthenticated()) {
		var avatarStyle = {
			background: 'url(' + user.get('avatar') + ')',
			'backgroundSize': 'cover'
		};

		var channelLink = 'https://www.stream.me/' + user.get('username');

		menu = (
			<ul>
				<li>
					<span className='welcome'>
						<a target='_blank' href={channelLink} title={user.get('username')}>
							<span style={avatarStyle} className='tinyAvatar' />
							{ user.get('username') }
						</a>
					</span>

					<a href='/logout' title='Sign out'>Sign Out</a>
				</li>
			</ul>
		);
	}

	React.render(
		<header>
			<nav>{menu}</nav>
			<h1>
				<a href='/' title='meep' className='meep-red'>
					<span className='logo' />
					<span className='title'>!meep</span>
				</a>
			</h1>
		</header>,
		el
	);
}
