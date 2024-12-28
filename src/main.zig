const std = @import("std");
const json = @import("json");

pub fn getConfigPath(allocator: std.mem.Allocator) ![]const u8 {
    const env_map = try allocator.create(std.process.EnvMap);
    env_map.* = try std.process.getEnvMap(allocator);

    var value = env_map.get("XDG_CONFIG_HOME") orelse "";
    if (value.len == 0) {
        const home = env_map.get("HOME") orelse ".";
        value = try std.fmt.allocPrint(allocator, "{s}/.config", .{home});
    }

    return try std.fmt.allocPrint(allocator, "{s}/tmux-sessionizer", .{value});
}

pub fn main() !void {
    var args = std.process.args();
    // skip program name
    _ = args.skip();

    const cmd = args.next() orelse return;
    const Case = enum { a, add, c, create, h, help, v, version, l, list, r, remove, d, delete };
    const case = std.meta.stringToEnum(Case, cmd) orelse return;

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    switch (case) {
        .a, .add, .c, .create => {
            const config = try readConfig(allocator);
            std.debug.print("Config:\n  session_name: {s}\n  session_path: {s}\n  windows: {s}\n", .{
                config.session_name,
                config.session_path,
                config.windows,
            });
        },
        .h, .help => help(),
        .v, .version => version(),
        .l, .list => list(),
        .r, .remove, .d, .delete => delete(),
    }
}

const Config = struct {
    session_name: []const u8 = "default",
    session_path: []const u8 = "test",
    windows: []const []const u8 = &[_][]const u8{"default"},
};

pub fn readConfig(allocator: std.mem.Allocator) !Config {
    var cwd = std.fs.cwd();
    cwd = try cwd.makeOpenPath(try getConfigPath(allocator), .{});

    // const file = try dir.openFile(sub_path: []const u8, flags: File.OpenFlags)
    const configFileName = "config.json";
    const open_flags = std.fs.File.OpenFlags{ .mode = .read_write };
    const file = blk: {
        const f = cwd.openFile(configFileName, open_flags) catch {
            // break :blk try cwd.createFile(configFileName, .{});
            const f = try cwd.createFile(configFileName, .{});

            // write in a default config
            const serialized = try json.toSlice(allocator, Config{});
            const fw = f.writer();
            _ = try fw.writeAll(serialized);

            break :blk f;
        };
        break :blk f;
    };
    defer file.close();

    const file_buffer = try file.readToEndAlloc(allocator, 1024);
    const deserialized = try json.fromSlice(allocator, Config, file_buffer);
    return deserialized.value;
}

pub fn create() !void {
    // Create directory if it doesn't exist
    std.debug.print("create\n", .{});
}

pub fn help() void {
    std.debug.print("help\n", .{});
}

pub fn version() void {
    std.debug.print("version\n", .{});
}

pub fn list() void {
    std.debug.print("list\n", .{});
}

pub fn delete() void {
    std.debug.print("delete\n", .{});
}

test "simple test" {}
