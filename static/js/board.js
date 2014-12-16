var app = angular.module('boardApp', ['ngResource','ngRoute','ui.bootstrap']);
app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
	$locationProvider.html5Mode(true);
	$locationProvider.hashPrefix('!');
	$routeProvider
	.when("/board", {
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
		redirectTo: "/board"
	});
}]);
app.factory("Employee", function($resource) {
	return $resource('/api/employees/:id');
});


app.controller('IndexController', ['$scope', '$modal', 'Employee', function($scope, $modal, Employee) {
	$scope.employees = Employee.query();
	$scope.newTransaction = function($scope) {
		var modal = $modal.open({
			templateUrl: '/static/partial/newTransaction.html',
		    	controller: 'TransactionController',
		    	size: 'lg',
	});
	};
}]);
app.controller('EmployeeController', ['$scope', 'Employee', function($scope, Employee) {
	$scope.employees = Employee.query();
	$scope.showNewEmployee = false;
	$scope.showAddEmployee = function() {
		$scope.showNewEmployee = true;  
	};
	$scope.addEmployee = function() {
		var e = new Employee();
		e.Name = $scope.newName;
		var saved = e.$save();
		saved.then(function() {
			$scope.showNewEmployee = false;
			$scope.newName = "";
			$scope.employees.push(e);
		}, function() {
			alert("new employee not saved");
		});
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
}]).directive('newFocus', function($timeout) {
        return function(scope, element, attrs) {
            scope.$watch('showNewEmployee', function(newValue) {
                console.log('changed');
                $timeout(function() {
                    newValue && element[0].focus();
                });
            }, true);
        };
});

app.factory("Location", function($resource) {
	return $resource('/api/locations/:location');
});

app.controller('LocationController', ['$scope', 'Location', function($scope, Location) {
	$scope.locations = Location.query();
	$scope.newLocation = new Location();
	$scope.addLocation = function() {
		var saved = $scope.newLocation.$save();
		saved.then(function() {
			$scope.showNewLocation = false;
			$scope.locations.push($scope.newLocation);
			$scope.newLocation = new Location();
		}, function() {
			alert("new location not saved");
		});
	};
}]);

app.controller('TransactionController', ['$scope', function($scope) {
	
}]);
