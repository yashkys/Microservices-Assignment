val scala3Version = "3.5.0"

lazy val root = project
  .in(file("."))
  .settings(
    name := "com.kys.bs2",
    version := "0.1.0-SNAPSHOT",

    scalaVersion := scala3Version,

    libraryDependencies += "org.scalameta" %% "munit" % "1.0.0" % Test
  )
Compile / PB.targets := Seq(
  scalapb.gen() -> (Compile / sourceManaged).value / "scalapb"
)
libraryDependencies += "com.mysql" % "mysql-connector-j" % "8.4.0"
libraryDependencies += "com.thesamet.scalapb" %% "scalapb-runtime" % scalapb.compiler.Version.scalapbVersion % "protobuf"
libraryDependencies += "com.thesamet.scalapb" %% "scalapb-runtime-grpc" % "0.11.15"
libraryDependencies += "io.grpc" % "grpc-netty" % "1.53.0"
libraryDependencies += "com.typesafe.scala-logging" %% "scala-logging" % "3.9.5"
libraryDependencies += "ch.qos.logback" % "logback-classic" % "1.4.11"



