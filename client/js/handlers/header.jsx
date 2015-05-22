var React = require('react'),
	helpers = require('./helpers');

module.exports = function(el) {
	React.render(
		<header>
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
