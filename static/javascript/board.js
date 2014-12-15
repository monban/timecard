var app = angular.module('boardApp', ['ngResource','ngRoute']);
app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
	$locationProvider.html5Mode(true);
	$locationProvider.hashPrefix('!');
	$routeProvider
	.when("/", {
		templateUrl: "/static/partial/board.html",
		controller: "IndexController"
	})
	.when("/employees", {
		templateUrl: "/static/partial/employees.html",
		controller: "EmployeeController"
	})
	.when("/locations", {
		templateUrl: "/static/partial/locations.html",
		controller: "LocationController"
	})
	.otherwise({
		redirectTo: "/"
	});
}]);
app.factory("Employee", function($resource) {
	return $resource('/api/employees/:id');
});

app.controller('IndexController', ['$scope', 'Employee', function($scope, Employee) {
	$scope.employees = Employee.query();
}]);
app.controller('EmployeeController', function($scope, Employee) {
	$scope.employees = Employee.query();
	$scope.addEmployee = function() {
		var e = new Employee();
		e.Name = $scope.newName;
		var saved = e.$save();
		saved.then(function() {
			
		}, function() {
			alert("new employee not saved");
		});
		$scope.employees.push(e);
	};

	$scope.deleteEmployee = function(employee) {
		var i = $scope.employees.indexOf(employee);
		var removed = Employee.delete({id: employee.id});
		removed.then(function() {
			$scope.employees.splice(i,1);
		}, function() {
			alert("problem");
		});
	}
});

app.controller('LocationController', ['$scope', function($scope) {
}]);

