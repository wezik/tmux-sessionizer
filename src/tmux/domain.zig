pub const TmuxSession = struct {
    name: []const u8,
    path: []const u8,
    env: TmuxEnvironment = TmuxEnvironment{},
    panes: []const TmuxPane = &[_]TmuxPane{TmuxPane{}},
};

pub const TmuxEnvironment = struct {
    cmds: []const []const u8 = &[_][]const u8{},
};

pub const TmuxPane = struct {
    name: []const u8 = "default",
    active: bool = false,
    env: TmuxEnvironment = TmuxEnvironment{},
};
