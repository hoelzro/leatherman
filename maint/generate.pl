#!/usr/bin/perl

use strict;
use warnings;
use autodie;

use File::Basename 'fileparse';
use File::Spec 'splitdir';
use Encode;
use JSON::PP;

no warnings 'uninitialized';

my %doc;

while (<STDIN>) {
   my $c = decode_json($_);

   my $doc_path = ($c->{path} =~ s/\.go$/.md/r);
   my $d = do { open my $fh, '<:encoding(UTF-8)', $doc_path; local $/; <$fh> };

   my ($tool, $dir) = fileparse($c->{path}, '.go');
   
   my ($cat) = ($dir =~ m{/internal/tool/([^/]+)/[^/]+});

   if (!$cat) {
      warn "no category for $tool, skipping...\n";
      next;
   }

   $doc{$cat}{$tool} = "$d";
}

open my $fh, '<:encoding(UTF-8)', 'maint/README_begin.md';
my $begin = "<!-- Code generated by maint/generate-README. DO NOT EDIT. -->\n" .
            do { local $/; <$fh> };
close $fh;

open $fh, '<:encoding(UTF-8)', 'maint/README_end.md';
my $end = do { local $/; <$fh> };
close $fh;

for my $category (keys %doc) {
   for my $tool (keys %{$doc{$category}}) {
      $doc{$category}{$tool} = "#### `$tool`\n\n$doc{$category}{$tool}\n"
   }
}

my %offsets;
my $offset = length $begin;
my $body = $begin;

for my $category (sort keys %doc) {
   $body .= "### $category\n\n";
   $offset += length(encode('UTF-8', "### $category\n\n", Encode::FB_CROAK));;

   for my $tool (sort keys %{$doc{$category}}) {
      $body .= $doc{$category}{$tool};
      my $length = length(encode('UTF-8', $doc{$category}{$tool}, Encode::FB_CROAK));
      $offsets{$tool} = "[$offset:" . ($offset + $length) . "]";
      $offset += $length;
   }
}

$body .= $end;

open my $readme, '>:encoding(UTF-8)', 'README.mdwn';
print $readme $body;

close $readme;

open my $help, '>:encoding(UTF-8)', 'help_generated.go';
$body =~ s/`/` + "`" + `/g;
print $help "package main\n\n" .
   "var commandReadme map[string][]byte\n" .
   "func init() {\n" .
   "\tcommandReadme = map[string][]byte{\n";

print $help qq(\t\t"$_": readme$offsets{$_},\n\n) for sort keys %offsets;

print $help "\t}\n";
print $help "}\n";

system 'go', 'fmt';
