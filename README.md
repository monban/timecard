timecard
========

A sign-in board / timecard application for casual office use

This is meant to be used on a touch-screen monitor at the entrance to your office. Employees can sign themselves
in and out with simple screen on-screen controls. It is not meant as a replacement for industrial timecard
punchers or any place with actual payroll impact. Rather it should be a place where you can simply glance and
see where your coworkers are at the moment and when they plan to return.

The backend is written in Golang and presents a simple JSON API. The front-end is HTML5 and Angular.js. By default
the Go backend will also serve the static front-end, future versions will allow you to disable this capability if
you want to distribute the statics through nginx, apache, or a cdn.
