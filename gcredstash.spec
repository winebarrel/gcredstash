%define  debug_package %{nil}

Name:		gcredstash
Version:	0.2.0
Release:	1%{?dist}
Summary:	gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

Group:		Development/Tools
License:	Apache License, Version 2.0
URL:		https://github.com/winebarrel/gcredstash
Source0:	%{name}_%{version}.tar.gz
# https://github.com/winebarrel/gcredstash/releases/download/v0.2.0/gcredstash_0.2.0.tar.gz

%description
gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

%prep
%setup -q -n %{name}

%build
make

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
install -m 755 gcredstash %{buildroot}/usr/bin/

%files
%defattr(-,root,root,-)
/usr/bin/gcredstash
