var app = angular.module('boardApp', ['ngResource']);
app.factory("Employee", function($resource) {
	return $resource('/api/employees/:id');
});
app.controller('employeesController', function($scope, Employee) {
	Employee.query(function(data) {
		$scope.employees = data;
	});

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
