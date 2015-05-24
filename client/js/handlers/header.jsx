var React = require('react'),
	helpers = require('./helpers'),
	user = require('../user');

module.exports = function(el) {
	var menu = '';

	if(user.isAuthenticated()) {
		menu = (<ul><li><a href='/logout' title='Sign out'>Sign Out</a></li></ul>);
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
