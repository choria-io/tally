%define debug_package %{nil}
%define pkgname {{cpkg_name}}
%define version {{cpkg_version}}
%define bindir {{cpkg_bindir}}
%define etcdir {{cpkg_etcdir}}
%define release {{cpkg_release}}
%define dist {{cpkg_dist}}
%define binary {{cpkg_binary}}
%define tarball {{cpkg_tarball}}

Name: %{pkgname}
Version: %{version}
Release: %{release}.%{dist}
Summary: The Choria Tally Service
License: Apache-2.0
URL: https://choria.io
Group: System Tools
Packager: R.I.Pienaar <rip@devco.net>
Source0: %{tarball}
BuildRoot: %{_tmppath}/%{pkgname}-%{version}-%{release}-root-%(%{__id_u} -n)
Requires(pre): /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%description
A service that passively listen for events on a Choria network and expose the
observed data via a Prometheus compatible exporter.

%prep
%setup -q

%build

%install
rm -rf %{buildroot}
%{__install} -d -m0755  %{buildroot}/etc/sysconfig
%{__install} -d -m0755  %{buildroot}/usr/lib/systemd/system
%{__install} -d -m0755  %{buildroot}/etc/logrotate.d
%{__install} -d -m0755  %{buildroot}%{bindir}
%{__install} -d -m0755  %{buildroot}%{etcdir}
%{__install} -d -m0755  %{buildroot}/var/log/%{pkgname}
%{__install} -m0644 dist/%{pkgname}.service %{buildroot}/usr/lib/systemd/system/%{pkgname}.service
%{__install} -m0644 dist/sysconfig %{buildroot}/etc/sysconfig/%{pkgname}
%{__install} -m0644 dist/%{pkgname}-logrotate %{buildroot}/etc/logrotate.d/%{pkgname}
%{__install} -m0755 %{binary} %{buildroot}%{bindir}/%{pkgname}
%{__install} -m0755 dist/config.json %{buildroot}%{etcdir}/config.json

touch %{buildroot}/var/log/%{pkgname}/%{pkgname}.log

%clean
rm -rf %{buildroot}

%post
if [ $1 -eq 1 ] ; then
  systemctl --no-reload preset %{pkgname} >/dev/null 2>&1 || :
fi

/bin/systemctl --system daemon-reload >/dev/null 2>&1 || :

if [ $1 -ge 1 ]; then
  /bin/systemctl try-restart %{pkgname} >/dev/null 2>&1 || :;
fi

%preun
if [ $1 -eq 0 ] ; then
  systemctl --no-reload disable --now %{pkgname} >/dev/null 2>&1 || :
  /usr/sbin/userdel %{pkgname} || :
fi

%pre
/usr/bin/getent group %{pkgname} || /usr/sbin/groupadd -r %{pkgname}
/usr/bin/getent passwd %{pkgname} || /usr/sbin/useradd -r -s /sbin/nologin -d /home/%{pkgname} -g %{pkgname} -c "Choria Tally Service" %{pkgname}

%files
%{bindir}/%{pkgname}
/etc/logrotate.d/%{pkgname}
/usr/lib/systemd/system/%{pkgname}.service
%attr(755, %{pkgname}, %{pkgname})/var/log/%{pkgname}
%config(noreplace) /etc/sysconfig/%{pkgname}
%config(noreplace) %{etcdir}

%changelog
* Fri Nov 08 2024 R.I.Pienaar <rip@devco.net>
- Initial Release
