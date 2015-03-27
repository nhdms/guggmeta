'use strict';

var app = angular.module('guggmeta', ['ngRoute', 'ngAria', 'ngAnimate', 'ui.bootstrap.pagination', 'angular-loading-bar']);

app.run(['$window', '$rootScope', function ($window, $rootScope) {
  $window.onload = function () {
    try {
      console.log("Welcome!\nYou can find the source code of guggmeta in https://github.com/sevein/guggmeta.\nPlease, send me your feedback!");
    } catch (e) {}
  };
  $rootScope.home = false;
  $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
    $rootScope.home = current.$$route.originalPath === '/';
  });
  $rootScope.top = function () {
    $window.scrollTo(0, 0);
  };
}]);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
  $routeProvider
    .when('/', {
      templateUrl: '/assets/partials/home.tmpl.html',
      controller: 'HomeCtrl'
    })
    .when('/analytics', {
      templateUrl: '/assets/partials/analytics-list.tmpl.html',
      controller: 'AnalyticsListCtrl'
    })
    .when('/submissions', {
      templateUrl: '/assets/partials/submission-list.tmpl.html',
      controller: 'SubmissionListCtrl',
      resolve: {
        response: ['$location', 'SubmissionService', function ($location, SubmissionService) {
          var query = '';
          var params = $location.search();
          if (params.hasOwnProperty('q')) {
            query = params.q;
          }
          return SubmissionService.search(query);
        }]
      }
    })
    .when('/submissions/:id', {
      templateUrl: '/assets/partials/submission-detail.tmpl.html',
      controller: 'SubmissionDetailCtrl',
      resolve: {
        response: ['$route', 'SubmissionService', function ($route, SubmissionService) {
          return SubmissionService.getOne($route.current.params.id);
        }]
      }
    })
    .otherwise('/');
  $locationProvider.html5Mode({
    enabled: true,
    requireBase: false
  });
}]);

app.controller('HeaderCtrl', ['$scope', '$location', '$window', function ($scope, $location, $window) {
  $scope.form = {};
  $scope.search = function () {
    if (angular.isUndefined($scope.query)) {
      return;
    }
    if (!$scope.query.length) {
      $location.path('/submissions');
      return;
    }
    if (/GH-\d/.test($scope.query)) {
      $location.path('/submissions/' + $scope.query);
      return;
    }
    $location.path('/submissions').search({ 'q': $scope.query });
  };
  $scope.$on('$routeChangeError', function () {
    $window.alert('Not found!');
    delete $scope.query;
  });
}]);

app.controller('HomeCtrl', ['$scope', function ($scope) {

}]);

app.controller('AnalyticsListCtrl', ['$scope', 'SubmissionService', function ($scope, SubmissionService) {
  SubmissionService.getAnalytics().then(function (response) {
    $scope.analytics = response.data;
  });
}]);

app.controller('SubmissionListCtrl', ['$scope', '$location', 'response', 'SubmissionService', function ($scope, $location, response, SubmissionService) {
  var params = $location.search();
  $scope.query = params.hasOwnProperty('q') ? params.q : undefined;
  var populate = function (resp) {
    $scope.submissions = resp.data.results;
    $scope.totalItems = resp.data.total;
  };
  populate(response);
  $scope.pageChanged = function () {
    SubmissionService.search($scope.query, $scope.currentPage).then(function (response) {
      populate(response);
      $scope.top();
    });
  };
}]);

app.controller('SubmissionDetailCtrl', ['$scope', '$routeParams', 'response', function ($scope, $routeParams, response) {
  $scope.id = $routeParams.id;
  $scope.submission = response.data;
  // Properties "file_name" and "author" are ignored or handled manually
  $scope.submissionFields = [
    { 'label': 'Creation date', 'property': 'creation_date' },
    { 'label': 'Creator', 'property': 'creator' },
    { 'label': 'Encrypted', 'property': 'encrypted' },
    { 'label': 'Size', 'property': 'file_size', 'suffix': ' bytes' },
    { 'label': 'Form', 'property': 'form' },
    { 'label': 'JavaScript', 'property': 'javascript' },
    { 'label': 'Keywords', 'property': 'keywords' },
    { 'label': 'Modification date', 'property': 'mod_date' },
    { 'label': 'Optimized', 'property': 'optimized' },
    { 'label': 'Page rot', 'property': 'page_rot' },
    { 'label': 'Page size', 'property': 'page_size' },
    { 'label': 'Pages', 'property': 'pages' },
    { 'label': 'PDF version', 'property': 'pdf_version' },
    { 'label': 'Producer', 'property': 'producer' },
    { 'label': 'Subject', 'property': 'subject' },
    { 'label': 'Suspects', 'property': 'suspects' },
    { 'label': 'Tagged', 'property': 'tagged' },
    { 'label': 'Title', 'property': 'title' },
    { 'label': 'UserProperties', 'property': 'user_properties' },
  ];
}]);

app.service('SubmissionService', ['$http', function ($http) {
  this.getOne = function (id) {
    return $http.get('/api/submissions/' + id + '/');
  };
  this.getAll = function () {
    return $http.get('/api/submissions/');
  };
  this.getAnalytics = function () {
    return $http.get('/api/submissions/analytics/');
  };
  this.search = function (query, page) {
    var params = {};
    if (angular.isDefined(query)) {
      params.q = query;
    }
    if (angular.isDefined(page)) {
      params.p = page;
    }
    return $http.get('/api/submissions/search/', {
      params: params
    });
  };
}]);
