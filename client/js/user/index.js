var User = function() {
	this.attributes = {
		authenticated: false
	};
};

User.prototype.set = function(attributes) {
	this.attributes = attributes;
	this.attributes.authenticated = true;
};

User.prototype.isAuthenticated = function() {
	return this.attributes.authenticated;
}


module.exports = new User();
