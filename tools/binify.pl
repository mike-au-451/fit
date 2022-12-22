#!/usr/bin/perl

=pod

Convert a file of bytes represented as hex digits to binary values.


=cut

use Modern::Perl;

while (my $line = <STDIN>) {
	$line = trim($line);
	next if length($line) == 0;
	next if substr($line, 0, 1) eq ';';

	my @elems = split(/\s+/, $line);

	foreach my $elem (@elems) {
		my $bin = hex($elem);
		print chr($bin);
	}
}

sub trim {
	my $str = shift // "";
	$str =~ s/^\s+//;
	$str =~ s/\s+$//;
	return $str;
}