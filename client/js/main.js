var page = require('page'),
	pubsub = require('./pubsub'),
	user = require('./user'),
	header = require('./handlers/header.jsx'),
	authentication = require('./handlers/authentication.jsx');

var pagechange = function(newPage) {
	return function(ctx, next) {
		pubsub.emit('page:change', {page: newPage});
		next();
	};
}

authentication(
	document.getElementsByTagName('menu')[0], function() {

	header(document.getElementsByTagName('header')[0]);

	page('/logout', function(){
		document.cookie = 'sessid=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
		window.location = '/';
	});

	if(user.isAuthenticated()) {
		page('/', pagechange('landing'), require('./handlers/home.jsx'));	
	} else {
		page('/', pagechange('landing'), require('./handlers/landing.jsx'));
	}
	
	// page('*')
	page();
});
