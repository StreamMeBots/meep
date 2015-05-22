var React = require('react'),
	helpers = require('./helpers');

module.exports = function(ctx) {
	React.render(
		<section className='jumbotron'>
			<div className='content'>
				<h1>Meet <span className='meep-red'>!meep</span></h1>
				<p>!meep is a bot for StreamMe to help channel owners greet and thank their viewers.</p>

				<div className='features'>
					<ul>
						<li>Greet your viewers and loyal viewers</li>
						<li>Tell viewers about your schedule when you are not streaming</li>
						<li>Set up commands to answer questions for your users</li>
					</ul>
				</div>

				<div className='actions'>
					<a className='button' href='/login' title='Sign in with StreamMe'>
						Sign in with StreamMe
					</a>
				</div>
			</div>
		</section>,
		helpers.getPageDiv()
	);
}
