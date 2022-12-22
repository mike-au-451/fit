#!/usr/bin/perl

use Modern::Perl;

foreach my $arg (@ARGV) {
	foreach my $ch (split //, $arg) {
		printf "%02x ", ord($ch);
	}
	print "\n";
}