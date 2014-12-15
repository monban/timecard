var app = angular.module('boardApp', ['ngResource','ngRoute']);
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

app.controller('IndexController', ['$scope', 'Employee', function($scope, Employee) {
	$scope.employees = Employee.query();
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

app.controller('LocationController', ['$scope', function($scope) {
}]);

