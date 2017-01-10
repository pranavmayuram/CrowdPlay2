var crowdPlay = angular.module('crowdPlay', []);

crowdPlay.run(function($rootScope) {
    $rootScope.baseName= "CrowdPlay";
    $rootScope.channelSocket = io();
    console.log("makin waves");
});

crowdPlay.controller('loginController', function($scope, $rootScope, $http, $window) {
    $scope.joinChannel = function(channelToJoin) {
        console.log("leggo Join");
        $rootScope.channelSocket.on('connect', function(response) {
           // Connected, let's sign-up for to receive messages for this room
           if (response) {
               console.log(response);
           }
           else {
               console.log("lel");
           }
           $rootScope.channelSocket.emit('joinChannel', channelToJoin);
        });
    };
    
    $scope.createChannel = function(channelToCreate) {
        console.log("leggo Create");
        $rootScope.channelSocket.on('connect', function(response) {
           // Connected, let's sign-up for to receive messages for this room
           if (response) {
               console.log(response);
           }
           else {
               console.log("lel");
           }
           $rootScope.channelSocket.emit('createChannel', channelToCreate);
        });
    };
});
