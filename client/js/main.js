var page = require('page'),
	pubsub = require('./pubsub'),
	header = require('./handlers/header.jsx'),
	authentication = require('./handlers/authentication.jsx');

var pagechange = function(newPage) {
	return function(ctx, next) {
		pubsub.emit('page:change', {page: newPage});
		next();
	};
}

header(document.getElementsByTagName('header')[0]);
authentication();

page('/', pagechange('landing'), require('./handlers/landing.jsx'));
// page('*')
page();
