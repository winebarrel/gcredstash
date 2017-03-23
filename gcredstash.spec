%define  debug_package %{nil}

Name:		gcredstash
Version:	0.3.0
Release:	1%{?dist}
Summary:	gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

Group:		Development/Tools
License:	Apache License, Version 2.0
URL:		https://github.com/winebarrel/gcredstash
Source0:	%{name}.tar.gz
# https://github.com/winebarrel/gcredstash/releases/download/v%{version}/gcredstash_%{version}.tar.gz

%description
gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

%prep
%setup -q -n src

%build
make

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/sbin
install -m 700 gcredstash %{buildroot}/usr/sbin/

%files
%defattr(700,root,root,-)
/usr/sbin/gcredstash
