pub const TmuxSession = struct {
    name: []const u8,
    path: []const u8,
    env: [][]const u8 = &[_][]const u8{},
    windows: []const TmuxWindow = &[_]TmuxWindow{TmuxWindow{}},
};

pub const TmuxWindow = struct {
    name: []const u8 = "default",
    active: bool = false,
    env: [][]const u8 = &[_][]const u8{},
};
